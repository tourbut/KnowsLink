const task = await taskSpace("connect-simple-qa-20261027");
const page = task.page("p1");
console.log(await page.url(), await page.title());
console.log((await page.snapshot()).slice(0, 600));
