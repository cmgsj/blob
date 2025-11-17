window.onload = function () {
  window.ui = SwaggerUIBundle({
    urls: [
      {
        name: "blob.api.v1.BlobService",
        url: "/docs/openapi/blob/api/v1/blob.openapi.json",
      },
    ],
    docExpansion: "list",
    defaultModelsExpandDepth: 0,
    defaultModelExpandDepth: 0,
    displayOperationId: false,
    displayRequestDuration: true,
    persistAuthorization: true,
    dom_id: "#swagger-ui",
    deepLinking: true,
    presets: [SwaggerUIBundle.presets.apis, SwaggerUIStandalonePreset],
    plugins: [SwaggerUIBundle.plugins.DownloadUrl],
    layout: "StandaloneLayout",
  });
};
