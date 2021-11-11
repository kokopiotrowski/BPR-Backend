/*
 * StockX - Bachelor Project
 *
 * This API is developed for a Bachelor Project Application named StockX.   - It both projects current stock market data and works as a learning tool while doing fake stock trading.    Authors of the project - Aleksander Stefan Bialik - Konrad Piotrowski  
 *
 * API version: 1.0.0
 * Contact: 280053@via.dk
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package swagger

import (
	"net/http"
)

func UserSettingsGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}
