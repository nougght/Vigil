import {NavLink} from "react-router-dom";
import styles from "./bottomNav.module.css";
import {useTheme} from "../hooks/useTheme.ts";
import type {Icon} from "@phosphor-icons/react";


export interface NavItem {
    id: string
    title: string
    iconSrc?: string
    icon?: Icon
    path: string
    countLabel?: number
}

export interface BottomNavData {
    iconSrc: string
    title: string
    items: NavItem[]
}


const BottomNavButton = ({data}: { data: NavItem }) => {
    return (
        <NavLink
            key={data.path}
            to={data.path}
            className={({isActive}) => `${styles.bottomNavButton} ${isActive ? styles.active : ""}`}
        >
            <div className={styles.imgCont}>
              {data.iconSrc && data.iconSrc != "" && <img src={data.iconSrc} width="15px" height="15px"/>}
              {data.icon && <data.icon size={16}/>}
            </div>
            <span className={styles.bottomNavButtonTitle}>{data.title}</span>
            {/*<span className={styles.bottomNavButtonCount}>{data.countLabel}</span>*/}
        </NavLink>
    )
}

export const BottomNav = ({data}: { data: BottomNavData }) => {
    const [_theme, _setTheme] = useTheme();  // TODO: move upper
    return (
        <div className={styles.bottomNavContainer}>
            <nav className={styles.bottomNav}>
                {
                    Array.from(data.items.values()).map((item) => (
                        <BottomNavButton
                            key={item.id}
                            data={item}
                        />
                    ))
                }
            </nav>
        </div>
    )
}
