import { NavLink } from "react-router-dom"
import styles from "./sideBar.module.css"
import { useCallback, useEffect, useRef, useState } from "react"
import { useTheme } from "../hooks/useTheme"


export interface NavItem {
    id: string
    title: string
    iconSrc?: string
    path: string
    countLabel?: number
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
            {data.iconSrc && data.iconSrc != "" && <img src={data.iconSrc} width="15px" height="15px" />}
            <span className={styles.sideBarButtonTitle}>{data.title}</span>
            <span className={styles.sideBarButtonCount}>{data.countLabel}</span>
        </NavLink>
    )
}

export const SideBar = ({ data }: { data: SideBarData }) => {
    const [width, setWidth] = useState<number | null>();
    const [widthVal, setWidthVal] = useState<string | null>();
    const [isResizing, setIsResizing] = useState(false);

    const [theme, setTheme] = useTheme();
    const barRef = useRef<HTMLBaseElement>(null);

    useEffect(() => {
        width != null && setWidthVal(`${width}px`)
    }, [width])

    useEffect(() => {
        if (isResizing) {
            document.body.style.cursor = 'w-resize';
        } else {
            document.body.style.cursor = 'default';
        }

        return () => {
            document.body.style.cursor = 'default';
        };
    }, [isResizing]);

    const handleTheme = () => {
        setTheme(theme == 'light' ? 'dark' : 'light');
    }

    const handleResize = useCallback((e: React.MouseEvent<HTMLDivElement>): void => {
        setIsResizing(true);
        e.preventDefault();
        const startWidth = barRef.current?.offsetWidth;
        // console.log("curr width = ", startWidth)
        const startX = e.clientX;
        if (barRef.current == null) {
            // console.log("ref null")
            return
        }

        const previousWidthStyle = barRef.current.style.width;
        barRef.current.style.width = 'min-content';
        const contentWidth = barRef.current.scrollWidth;
        barRef.current.style.width = previousWidthStyle;

        const doResize = (e: MouseEvent) => {
            const newWidth = startWidth! + (e.clientX - startX);

            // console.log("newWidth = ", newWidth)
            const minAllowedWidth = contentWidth - 20;

            if (newWidth >= minAllowedWidth && newWidth < 800) {
                barRef.current!.style.width = `${newWidth}px`;
            }
        };

        const stopResize = () => {
            window.removeEventListener('mousemove', doResize);
            setIsResizing(false);
            window.removeEventListener('mouseup', stopResize);
            setWidth(barRef.current?.offsetWidth)
        };
        window.addEventListener('mousemove', doResize);
        window.addEventListener('mouseup', stopResize);

    }, [width])

    return (
        <aside className={styles.sideBar}
            ref={barRef}
            style={{
                '--width-val': widthVal,
            } as React.CSSProperties}>
            <div className={styles.barContainer}>
                <nav >
                    {
                        Array.from(data.items.values()).map((item) => (
                            <SideBarButton
                                key={item.id}
                                data={item}
                            />
                        ))
                    }
                </nav>
                <div className={styles.bottomContainer}>
                    <SideBarButton
                        key=""
                        data={{
                            id: "",
                            title: "Настройки",
                            path: "/settings"
                        }}
                    />
                    <button className={styles.themeButton}
                        onClick={handleTheme}
                    >
                        Поменять тему
                    </button>
                </div>
            </div>
            <div className={styles.resizeArea}
                onMouseDown={handleResize}
            />
        </aside>
    )
}
