function colorBlink(item, timeout, idleColor, blinkColor) {
  item.classList.replace(idleColor, blinkColor);
  setTimeout(() => {
    item.classList.replace(blinkColor, idleColor);
  }, timeout);
}
