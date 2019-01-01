const SERVER_HOST = "http://localhost:3333"
const createFetch = (route, method) => {
 return (data) => {
  const options = {
    method,
    headers: {
      Accept: 'application/json',
      'Content-Type': 'application/json',
    },
    ...(data !== undefined ? { body: JSON.stringify(data) } : {}),
  }
  return fetch(`${SERVER_HOST}${route}`, options).then(res => res.json())
 }
}

export const createNewSession = createFetch('/api/session', 'POST')

export const getSession = async (id, adminID) => {
  const adminStr = adminID ? `adminID=${adminID}` : ''
  return createFetch(`/api/session?id=${id}&${adminStr}`, 'GET')()
}


//   {
//     "dateCreated": "2018-12-31 06:23:47.193119 +0000 UTC",
//     "ID": "b8b1b9a2-1bb7-4b7f-8ebb-276e0c7e2aa9",
//     "storyPoints": [
//       1,
//       2,
//       3
//     ],
//     "issues": [
//       {
//         "issueTitle": "",
//         "issueID": "db1f7fd4-aff3-4495-80fa-ef842c7eda71",
//         "estimations": {
//           "13": 0,
//           "123": 1231,
//           "": 123,
//           "bar": 456,
//           "foo": 123
//         }
//       }
//     ],
//     "selectedIssue": "db1f7fd4-aff3-4495-80fa-ef842c7eda71"
//   },

export const updateSession = async (id, adminID, newSession) => {
  const adminStr = adminID ? `adminID=${adminID}` : ''
  return createFetch(`/api/session?id=${id}&${adminStr}`, 'PUT')(JSON.stringify(newSession))
}
