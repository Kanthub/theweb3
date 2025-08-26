try {
  setTimeout(function () {
    throw new Error("异步出错");
  }, 0);
} catch (error) {
  console.log("error");
}
console.log("out try catch");
