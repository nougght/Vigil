import { useEffect, useRef, useState } from "react";
import { Cell, Legend, Pie, PieChart, ResponsiveContainer, Sector, Tooltip } from "recharts";


interface PieRow {
    Name: string,
    Value: number,
    fill: string,
}

export default function SimplePieChart({ Title, PieData }: { Title: string, PieData: PieRow[] }) {

    const first = useRef(true);
    useEffect(() => { first.current = false; }, []);

    return (
        <div style={{ width: '100%', height: 300 }}>
            <h3>{Title}</h3>
            <ResponsiveContainer>
                <PieChart>
                    <Pie
                        data={PieData}
                        cx="50%" // Center X
                        cy="50%" // Center Y
                        innerRadius={0} // Set greater than 0 for a Donut Chart
                        outerRadius={80}
                        dataKey="Value"
                        nameKey="Name"
                        // label={({ name }) => `${name}`}
                        label={({ name, percent }) => percent == null || percent == 0 ? "" :`${name} ${(percent * 100).toFixed(0)}%`}
                        // shape={renderShape}
                        animationDuration={300}
                        isAnimationActive={first.current}
                    >
                    </Pie>
                    <Tooltip />
                    <Legend
                        // payload={PieData.map(r => ({
                        //     value: r.Name,
                        //     type: 'square',
                        //     color: r.fill,
                        //     inactive: r.Value === 0,
                        // }))}
                    />
                </PieChart>
            </ResponsiveContainer>
        </div>
    );
}

