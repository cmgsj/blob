window.onload = function () {
  window.ui = SwaggerUIBundle({
    urls: [{ name: "blob.v1", url: "/docs/schemas/blob.v1.swagger.json" }],
    dom_id: "#swagger-ui",
    deepLinking: true,
    presets: [SwaggerUIBundle.presets.apis, SwaggerUIStandalonePreset],
    plugins: [SwaggerUIBundle.plugins.DownloadUrl],
    layout: "StandaloneLayout",
  });
};
