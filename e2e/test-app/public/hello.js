console.log("Hello, world!");

(function () {
  window.setTimeout(() => {
    const data = fetch("/api/data", { method: "POST" })
      .then((res) => res.json())
      .then((data) => {
        console.log("Fetched data:", data);
      });
  }, 2000);
})();
