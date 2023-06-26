async function main(systemData) {
  console.log("<---", systemData);
  return "exec ok";
}
process.on("message", async (systemData) => {
  const result = await main(systemData);
  process.send(result);
  // process.exit();
});

