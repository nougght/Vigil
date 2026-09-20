import { useEffect, useState } from "react"


export const useTheme = (): [string, React.Dispatch<React.SetStateAction<string>>] => {
    const [theme, setTheme] = useState('light')
    useEffect(() => {
        document.documentElement.setAttribute('app-theme', theme);
        document.documentElement.style.colorScheme = theme;
    });

    return [theme, setTheme];
}