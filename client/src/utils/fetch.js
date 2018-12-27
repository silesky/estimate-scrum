import { toBool } from '../utils'
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
