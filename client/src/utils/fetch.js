const SERVER_HOST = 'http://localhost:3333';
const createFetch = (route, method) => {
  return data => {
    const options = {
      method,
      headers: {
        'Content-Type': 'application/json',
      },
      ...(data !== undefined ? { body: JSON.stringify(data) } : {}),
    };
    console.log('options', options)
    return fetch(`${SERVER_HOST}${route}`, options).then(res => {
      if (!res.ok) throw res;
      return res.json();
    });
  };
};

export const createNewSession = createFetch('/api/session', 'POST');

export const getSession = async (id, adminID) => {
  const adminStr = adminID ? `adminID=${adminID}` : '';
  return createFetch(`/api/session?id=${id}&${adminStr}`, 'GET')();
};
export const addEstimation = data => {
 return createFetch(`/api/estimation`, 'POST')(data)
}
//
// const validateEstimation = () => {
//   "username": "JIM",
// 	"sessionID": "fce5fbc2-ce78-4da4-a7da-6e3ddf678571",
// 	"issueID": "b98393de-3d75-46f0-a439-3a6a2522891a",
// 	"estimationValue": 123
// }

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
  const adminStr = adminID ? `adminID=${adminID}` : '';
  return createFetch(`/api/session?id=${id}&${adminStr}`, 'PUT')(newSession);
};
