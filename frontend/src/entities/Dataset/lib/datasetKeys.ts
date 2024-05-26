class DatasetKeys {
  root = ['dataset']

  list = () => [...this.root, 'list']
}

export const datasetKeys = new DatasetKeys()
