import { getGradientColor } from '@monitoring-system/shared/src/util/gradientColor'
import styles from './loadBar.module.css'

export const LoadBar = ({ value, width, height }: { value: number, width?: number, height?: number }) => {
    const wdValue =width == undefined ? '100%' : `${width}px`
    const hgValue = height == undefined ? '100%' : `${height}px`

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