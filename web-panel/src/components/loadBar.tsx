import { getGradientColor } from '@monitoring-system/shared/src/util/gradientColor'
import styles from './loadBar.module.css'
import { useEffect, useState } from 'react'

export const LoadBar = ({ value, width, height }: { value: number, width?: number, height?: number }) => {
    const [wdValue, setWd] = useState<string>('100%')
    const [hgValue, setHg] = useState<string>('20px')
    useEffect(() => {
        setWd(width == undefined ? '100%' : `${width}px`)
    }, [width])
    useEffect(() => {
        setHg(height == undefined ? '100%' : `${height}px`)
    }, [height])

    return (
        <div className={styles.container}
            style={{
                '--width': wdValue,
                '--height': hgValue,
            } as React.CSSProperties}>
            <div className={styles.filler}
                style={{
                    '--load-percent': `${value}%`,
                    '--color': `${getGradientColor(["#62e09b", "#e0be62", "#e06262"], value)}`,
                } as React.CSSProperties} />
        </div>
    )
}