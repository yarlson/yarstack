addEventListener("DOMContentLoaded", () => {
  const dot = document.createElement("div");
  dot.style.cssText =
    "position:fixed;left:-40px;top:-40px;width:20px;height:20px;margin:-10px 0 0 -10px;border-radius:50%;" +
    "background:rgba(255,72,72,.55);border:2px solid #fff;box-shadow:0 0 6px rgba(0,0,0,.6);" +
    "pointer-events:none;z-index:2147483647;transition:transform .15s";
  document.documentElement.append(dot);
  addEventListener("mousemove", (event) => {
    dot.style.display = "";
    dot.style.left = `${event.clientX}px`;
    dot.style.top = `${event.clientY}px`;
  }, true);
  addEventListener("mousedown", () => { dot.style.transform = "scale(1.8)"; }, true);
  addEventListener("mouseup", () => { dot.style.transform = ""; }, true);
  document.addEventListener("mouseout", (event) => {
    if (!event.relatedTarget) dot.style.display = "none";
  }, true);
});
