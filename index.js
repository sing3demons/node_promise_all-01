const population = {
  th: [
    { id: 1, name: 'nameth1' },
    { id: 2, name: 'nameth2' },
    { id: 3, name: 'nameth3' },
    { id: 4, name: 'nameth4' },
  ],
  en: [
    { id: 5, name: 'enname1' },
    { id: 6, name: 'enname2' },
    { id: 7, name: 'enname3' },
    { id: 8, name: 'enname4' },
  ],
  my: [
    { id: 9, name: 'myname1' },
    { id: 10, name: 'myname2' },
    { id: 11, name: 'myname3' },
    { id: 12, name: 'myname4' },
  ],
  km: [
    { id: 13, name: 'kmname1' },
    { id: 14, name: 'kmname2' },
    { id: 15, name: 'kmname3' },
    { id: 16, name: 'kmname4' },
  ],
}

async function createData(body) {
  return new Promise((resolve, reject) => {
    setTimeout(() => {
      resolve({ msg: 'success', data: body })
    }, 1000)
  })
}

async function main() {
  const response = []

  for (const key in population) {
    if (population.hasOwnProperty(key)) {
      for (const body of population[key]) {
        response.push(createData(body))
      }
    }
  }

  if (response.length) {
    const data = await Promise.all(response)
    console.log(data)
  }
}
main()
