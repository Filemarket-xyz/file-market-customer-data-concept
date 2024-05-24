import { ComponentType, createContext } from 'react'

export interface AppModalProps {
  open: boolean
  onClose: () => void
}

export interface ModalInstance<Props extends AppModalProps = AppModalProps> {
  Component: ComponentType<Props>
  open: boolean
  props?: Omit<Props, 'open' | 'onClose'> | undefined
}

export type OpenModalHandler = <Props extends AppModalProps>(instance: Omit<ModalInstance<Props>, 'id' | 'open'>) => void

export interface ModalContext {
  openModal: OpenModalHandler
}

export const ModalContext = createContext<ModalContext>({
  openModal: () => {},
})
