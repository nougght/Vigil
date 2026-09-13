import { NavLink } from "react-router-dom"
import styles from "./sideBar.module.css"


export interface NavItem {
    id: string
    title: string
    iconSrc: string
    path: string
    countLabel: number
}

export interface SideBarData {
    iconSrc: string
    title: string
    items: NavItem[]
    // TODO: add more sections
}

// interface SideBarProps {
//     data: SideBarData
//     onSwitch: (id: string) => void
//     activeId: string
// }


const SideBarButton = ({ data }: { data: NavItem }) => {
    return (
        <NavLink
            key={data.path}
            to={data.path}
            className={({ isActive }) => `${styles.sideBarButton} ${isActive ? styles.active : ""}`}
        >
            {data.iconSrc != "" && <img src={data.iconSrc} width="15px" height="15px" />}
            <span className={styles.sideBarButtonTitle}>{data.title}</span>
            <span className={styles.sideBarButtonCount}>{data.countLabel}</span>
        </NavLink>
    )
}

export const SideBar = ({ data}: {data: SideBarData}) => {
    return (
        <aside className={styles.sideBar}>
            <nav>
                {
                    Array.from(data.items.values()).map((item) => (
                        <SideBarButton
                            key={item.id}
                            data={item}
                        />
                    ))
                }
            </nav>
        </aside>
    )
}
