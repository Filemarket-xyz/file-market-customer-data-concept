class UserKeys {
  root = ['user']

  connectMessage = (address: string) => [...this.root, 'connectMessage', address]

  loginBySignature = (address: string) => [...this.root, 'loginBySignature', address]

  logout = () => [...this.root, 'logout']
}

export const userKeys = new UserKeys()
