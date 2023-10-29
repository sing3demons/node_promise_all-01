import { connect, getCollection } from './db.js'
import express from 'express'
import NanoIdService from './nanoid.js'

const app = express()
const nano = new NanoIdService()

async function main() {
  await connect()

  app.use(express.json())

  app.get('/', async (req, res) => {
    const result = await getCollection().find({}).toArray()
    res.status(200).json(result)
  })

  app.post('/example', async (req, res) => {
    const { name } = req.body
    console.log(req.body)
    const id = nano.randomNanoId()
    const result = await getCollection().insertOne({
      id,
      name,
      createDate: new Date(new Date().toUTCString()),
      updateDate: new Date(new Date().toUTCString()),
    })
    result.insertedId
    const data = await getCollection().findOne({ _id: result.insertedId })
    res.status(201).json(data)
  })

  app.delete('/example/:id', async (req, res) => {
    const { id } = req.params

    const result = await getCollection().findOneAndUpdate(
      { id },
      { $set: { deleteDate: new Date(new Date().toUTCString()) } },
      { upsert: true, returnDocument: 'after' }
    )

    res.status(200).json(result)
  })

  app.listen(3000, () => console.log('Server listening on port 3000'))
}
main()
