import * as jdenticon from 'jdenticon'
import React, { useEffect, useRef } from 'react'

export interface AddressIconProps {
  address: string
  size?: number
}

export const AddressIcon: React.FC<AddressIconProps> = ({ address, size = 56 }) => {
  const ref = useRef<SVGSVGElement>(null)

  useEffect(() => {
    // @ts-expect-error bad types in 'jdenticon'
    jdenticon.update(ref.current, address)
  }, [address])

  return (
    <svg
      ref={ref}
      data-jdenticon-value={address}
      height={size}
      width={size}
    />
  )
}
