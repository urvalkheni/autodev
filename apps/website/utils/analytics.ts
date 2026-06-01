/**
 * Zero-dependency telemetry tracking client for Autodevs.dev website.
 * Pings counterapi.dev to record visitor activation funnels.
 */

export const trackInstall = (method: string) => {
  // Fire main install metric
  fetch("https://api.counterapi.dev/v1/heetmehta18-autodev/installs/up").catch(
    (err) => console.error("Analytics failure:", err),
  );

  // Fire specific method metric
  fetch(
    `https://api.counterapi.dev/v1/heetmehta18-autodev/install_${method.toLowerCase()}/up`,
  ).catch(() => {});
};

export const trackTemplateView = (templateName: string) => {
  fetch(
    `https://api.counterapi.dev/v1/heetmehta18-autodev/template_${templateName.toLowerCase()}_views/up`,
  ).catch(() => {});
};

export const trackTemplateCopy = (templateName: string) => {
  fetch(
    `https://api.counterapi.dev/v1/heetmehta18-autodev/template_${templateName.toLowerCase()}_copies/up`,
  ).catch(() => {});
};
