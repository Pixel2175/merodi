-- HTTP are available only with `merodi serve`
-- Use routes to handle requests and return responses

merodi.http.route("/path", function()
    return {
        mime = "text/plain",
        value = "response",
        status = 200
    }
end)

-- Registers a handler for a URL path.
-- The handler returns a response table:
--   mime   (string) - response Content-Type
--   value  (string) - response body
--   status (number) - HTTP status code
-- Defaults: mime = "text/plain", status = 200.

merodi.http.method()
-- Returns the current request method, such as "GET" or "POST".
-- Only available inside a route.

merodi.http.request.body()
-- Returns the request body as a string.

merodi.http.request.remote_addr()
-- Returns the client's IP address.

-- Server settings
merodi.http.settings.host.set("localhost")
-- Sets the host the server listens on.

merodi.http.settings.port.set("8866")
-- Sets the port the server listens on.

merodi.http.settings.host.get()
-- Returns the configured host.

merodi.http.settings.port.get()
-- Returns the configured port.
