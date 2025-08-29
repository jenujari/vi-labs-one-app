package router

import (
	"fmt"
	"log"
	"net/http"
	"server/config"
	"server/helpers"
	"server/repository"
	"server/sqlc"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
	kiteconnect "github.com/zerodha/gokiteconnect/v4"
	"golang.org/x/sync/errgroup"
)

func runRsiObvFilterForDay(w http.ResponseWriter, r *http.Request) {
	er := onetimeSevenFiftySync()
	if er != nil {
		http.Error(w, "Failed to sync 750 symbols: "+er.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	// repo := repository.GetRepository()
	// kc := helpers.GetKiteClient()

	// // conf := config.GetConfig()

	// accessToken, err := repo.GetZerodhaAccessToken()
	// if err != nil {
	// 	http.Error(w, "Failed to get token : "+err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// kc.SetAccessToken(accessToken)

	// sevenFifty, err := repo.GetSevenFiftySymbols()
	// if err != nil {
	// 	http.Error(w, "Failed to get 750 symbols from DB: "+err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// instruments, err := kc.GetInstrumentsByExchange("NSE")

	// if err != nil {
	// 	http.Error(w, "Failed to get instruments : "+err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// for _, sf := range sevenFifty {
	// 	for _, inst := range instruments {
	// 		if sf.Symbol == inst.Tradingsymbol && inst.InstrumentType == "EQ" {
	// 			ur := new(sqlc.UpdateSymbolsParams)
	// 			ur.Symbol = sf.Symbol
	// 			ur.FullName = pgtype.Text{String: inst.Name, Valid: true}
	// 			ur.InstrumentToken = pgtype.Int4{Int32: int32(inst.InstrumentToken), Valid: true}
	// 			ur.ExchangeToken = pgtype.Int4{Int32: int32(inst.ExchangeToken), Valid: true}

	// 			err := repo.UpdateSymbols(ur)
	// 			if err != nil {
	// 				http.Error(w, "Failed to update symbols : "+err.Error(), http.StatusInternalServerError)
	// 				return
	// 			}
	// 		}
	// 	}
	// }

	// // helpers.DumpHTML(w, instruments)

	// w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(instruments)
	// w.WriteHeader(http.StatusOK)
}

func refreshZerodhaToken(w http.ResponseWriter, r *http.Request) {
	repo := repository.GetRepository()
	kc := helpers.GetKiteClient()
	conf := config.GetConfig()

	token, err := repo.GetZerodhaAccessToken()

	if err == nil && len(token) > 0 {
		w.WriteHeader(http.StatusOK)
		return
	}

	requestToken, err := helpers.GetRequestTokenUsingBrowser()
	if err != nil {
		http.Error(w, "Failed to get auth token : "+err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := kc.GenerateSession(requestToken, conf.Secret.ApiSecret)
	if err != nil {
		http.Error(w, "Failed to authenticate: "+err.Error(), http.StatusInternalServerError)
		return
	}

	kc.SetAccessToken(data.AccessToken)
	log.Println("data.AccessToken", data.AccessToken)

	repo.SaveZerodhaAccessToken(data.AccessToken)

	_, err = kc.GetUserMargins()
	if err != nil {
		http.Error(w, "Failed to get user margins : "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func checkZerodhaTokenValidity(w http.ResponseWriter, r *http.Request) {
	repo := repository.GetRepository()
	accessToken, err := repo.GetZerodhaAccessToken()

	if err != nil {
		http.Error(w, "Failed to get token : "+err.Error(), http.StatusInternalServerError)
		return
	}

	kc := helpers.GetKiteClient()

	kc.SetAccessToken(accessToken)

	_, err = kc.GetUserMargins()
	if err != nil {
		http.Error(w, "Failed to get user margins : "+err.Error(), http.StatusInternalServerError)
		return
	}

	// config.GetLogger().Println(resp)

	w.WriteHeader(http.StatusOK)
}

func onetimeSevenFiftySync() error {
	repo := repository.GetRepository()
	kc := helpers.GetKiteClient()

	// conf := config.GetConfig()

	accessToken, err := repo.GetZerodhaAccessToken()
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}

	kc.SetAccessToken(accessToken)

	sevenFifty, err := repo.GetSevenFiftySymbols()
	if err != nil {
		return fmt.Errorf("failed to get 750 symbols from DB: %w", err)
	}

	instruments, err := kc.GetInstrumentsByExchange("NSE")

	if err != nil {
		return fmt.Errorf("failed to get instruments: %w", err)
	}

	g := new(errgroup.Group)

	for _, sf := range sevenFifty {
		sf := sf
		if !sf.InstrumentToken.Valid {
			config.GetLogger().Println("Updating symbol:", sf.Symbol, " ", sf.InstrumentToken.Valid)
			g.Go(func() error {
				return updateSymbol(sf, instruments, repo)
			})
		}
	}

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}

func updateSymbol(sf sqlc.TblSevenFifty, instruments kiteconnect.Instruments, repo *repository.Repository) error {
	for _, inst := range instruments {
		if strings.HasPrefix(inst.Tradingsymbol, sf.Symbol) && inst.InstrumentType == "EQ" {
			ur := new(sqlc.UpdateSymbolsParams)
			ur.Symbol = sf.Symbol
			ur.FullName = pgtype.Text{String: inst.Name, Valid: true}
			ur.InstrumentToken = pgtype.Int4{Int32: int32(inst.InstrumentToken), Valid: true}
			ur.ExchangeToken = pgtype.Int4{Int32: int32(inst.ExchangeToken), Valid: true}

			err := repo.UpdateSymbols(ur)
			if err != nil {
				return fmt.Errorf("failed to update symbol %s: %w", sf.Symbol, err)
			}
		}
	}
	return nil
}
