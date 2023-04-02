const express = require('express')
const app = express()
const port = 18009
app.use(express.json());
app.post("/testt",(req,res)=>{
  console.log(req.body)
  console.log(req.headers)
  res.write("1234")
  res.end()
})
app.post("/relay-admin-panel/lpnode_admin_panel/register_lp",(req,res)=>{
  console.log(req.body)
  setTimeout(() => {
    res.write(JSON.stringify({
      code:0,
      message:"ok",
      lp_id_fake:"100",
      name:"lpname",
      profile:"LP profile",
      lpnode_api_key:"eiirrooror8",
      relay_api_key:"jsdjkfdjjdfs112"
    }))
    res.end()
  }, 500);
})
app.post("/lpnode/lpnode_admin_panel/relayAccount/register",(req,res)=>{
  console.log(req.body)
  setTimeout(() => {
    res.write(JSON.stringify({
      code:0,
      message:"ok",
      lp_id_fake:"100",
      name:"lpname",
      profile:"LP profile",
      lpnode_api_key:"eiirrooror8",
      relay_api_key:"jsdjkfdjjdfs112"
    }))
    res.end()
  }, 500);

})

app.post("/lpnode/lpnode_admin_panel/configLp",(req,res)=>{
  console.log(req.body)
  setTimeout(() => {
    res.write(JSON.stringify({
      code:0,
      message:"ok",
    }))
    res.end()
  }, 500);

})

app.post("/lpnode_admin_panel/set_wallet",(req,res)=>{
  console.log(req.body)
  setTimeout(() => {
    res.write(JSON.stringify({
      code:0,
      message:"ok",
    }))
    res.end()
  }, 500);

})
app.listen(port)

