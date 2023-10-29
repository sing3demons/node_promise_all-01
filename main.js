import axios from 'axios'

async function createData(body) {
  try {
    const { data } = await axios.post('http://localhost:3000/example', body, {
      headers: {
        'Content-Type': 'application/json',
      },
    })
    return data
  } catch (error) {
    console.log(error.message)
  }
}

async function createAsyncData(data) {
  try {
    const response = []

    for (let i = 0; i < data.length; i++) {
      const body = data[i]
      response.push(createData(body))
    }

    if (response.length) {
      const resultData = await Promise.all(response)
      resultData.forEach((data) => {
        console.log(data)
      })
    }
  } catch (error) {
    console.log(error.message)
  }
}

async function main() {
  const data = []
  for (let i = 0; i < 900; i++) {
    data.push({ name: `name_${i}` })
  }

  const start = performance.now()
  await mainCreateAsyncData(data) // Time: 14077.43397796154
  // await createAsyncData(data) // Time: 4774.951998949051

  const end = performance.now()
  console.log(`Time: ${end - start}`)
}
main()

async function mainCreateAsyncData(data) {
  for (let i = 0; i < data.length; i++) {
    const body = data[i]
    const result = await createData(body)
    console.log(result)
  }
}
