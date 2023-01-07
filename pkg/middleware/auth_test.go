package middleware_test

// TODO
// func TestAuth(t *testing.T) {
// 	conf, err := config.LoadFile("dev")
// 	require.NoError(t, err)
// 	h := app.NewHandler(conf)
// 	h.Routes()
// 	app.SetupMiddleware(h)
// 	defer h.Cleanup()

// 	// Create request
// 	u := url.URL{
// 		Path: "/api",
// 	}
// 	q := u.Query()
// 	q.Set("token", "123")
// 	u.RawQuery = q.Encode()
// 	req, err := http.NewRequest("GET", u.String(), nil)
// 	require.NoError(t, err)

// 	// Record handler response
// 	rec := httptest.NewRecorder()
// 	h.HTTPHandler.ServeHTTP(rec, req)

// 	// Verify response
// 	require.Equal(t, rec.Code, http.StatusOK)
// }
