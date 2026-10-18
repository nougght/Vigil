import  styles from "./TabBar.module.css"

export interface TabBarProps {
    className?: string;
    tabs: string[];
    onSwitch: (index: number) => void;
    activeTab: number;
}

interface TabButtonProps {
    text: string;
    onClick: () => void;
    isActive: boolean;
}

const TabButton = ({text, onClick, isActive}: TabButtonProps) => {
    return (
        <button className={styles.TabButton + (isActive ? ` ${styles.active}` : '')} onClick={onClick}>
            {text}
        </button>
    )
}


const TabBar = ({ className, tabs, onSwitch, activeTab }: TabBarProps) => {
    return (
        <div className={`${styles.TabBar} ${className ?? ""}`}>
            {tabs.map((tab, i) => (
                <TabButton key={tab} text={tab} onClick={() => onSwitch(i)} isActive={i === activeTab} />
            ))}
        </div>
    )
}

export default TabBar;