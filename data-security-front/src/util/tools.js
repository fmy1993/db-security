import crypto from 'crypto'
import store from '../store/index'
import JSEncrypt from 'jsencrypt'

export function getHash (str) {
  const md5 = crypto.createHash('md5')
  md5.update(str)
  return md5.digest('hex')
}

export function getRsa (str, publicKey) {
  const jse = new JSEncrypt()
  jse.setPublicKey(publicKey)
  return jse.encrypt(str)
}

export function refreshToken (response) {
  if (response.data.token != null) {
    store.commit('refreshToken', response.data.token)
  }
}
