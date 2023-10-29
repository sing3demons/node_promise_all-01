import { Collection, MongoClient } from 'mongodb'

const url = 'mongodb://mongo1:27017,mongo2:27018,mongo3:27019/mydatabase?replicaSet=my-replica-set'

const client = new MongoClient(url)

interface Example {
  id: string
  name: string
}

export async function connect() {
  try {
    await client.connect()
    console.log('Connected successfully to the database')
  } catch (error) {
    console.error('Failed to connect to the database')
    console.error(error)
  }
}

export function getCollection(): Collection<Example> {
  return client.db('microservice_db').collection<Example>('exampleDb')
}
