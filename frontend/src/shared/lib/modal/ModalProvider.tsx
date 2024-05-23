import React, { PropsWithChildren, useCallback, useRef, useState } from 'react'

import { ModalContext, ModalInstance, OpenModalHandler } from './ModalContext'

export const ModalProvider: React.FC<PropsWithChildren> = ({ children }) => {
  const modalId = useRef(0)
  const [instances, setInstances] = useState<Map<number, ModalInstance<any>>>(new Map())

  const openModal: OpenModalHandler = useCallback(({ Component, props }) => {
    setInstances((instances) => {
      return new Map(instances).set(modalId.current++, {
        Component,
        props,
        open: true,
      })
    },
    )
  }, [])

  const closeModal = (id: number) => {
    // closing modal with animation
    setInstances((instances) => {
      const instance = instances?.get(id)
      if (!instance) return instances

      return new Map(instances).set(id, {
        ...instance,
        open: false,
      })
    })

    // deleting modal after animation
    setTimeout(() => {
      setInstances((instances) => {
        instances?.delete(id)

        return new Map(instances)
      })
    }, 1000)
  }

  return (
    <ModalContext.Provider value={{ openModal }}>
      {children}
      {Array.from(instances.entries()).map(([id, { Component, props, open }]) => (
        <Component
          key={id}
          {...props}
          open={open}
          onClose={() => closeModal(id)}
        />
      ))}
    </ModalContext.Provider>
  )
}
