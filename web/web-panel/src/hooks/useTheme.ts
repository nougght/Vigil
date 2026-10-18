import { useEffect, useState } from "react"


const media = window.matchMedia('(prefers-color-scheme: dark)')

// TODO: add saving
export const useTheme = (): [string, React.Dispatch<React.SetStateAction<string>>] => {
    const [theme, setTheme] = useState(media.matches ? 'dark' : 'light')

    const onChange = () => setTheme(
        media.matches ? 'dark' : 'light'
    )

    useEffect(() => {
        document.documentElement.setAttribute('app-theme', theme);
        document.documentElement.style.colorScheme = theme;
        media.addEventListener('change', onChange)
        return () => media.removeEventListener('change', onChange)
    }, [theme]);

    return [theme, setTheme];
}