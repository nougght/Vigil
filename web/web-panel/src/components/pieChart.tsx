import { useEffect, useRef,  } from "react";
import { Legend, Pie, PieChart, ResponsiveContainer, Tooltip, type PieLabelRenderProps } from "recharts";


interface PieRow {
    Name: string,
    Value: number,
    fill: string,
}

export default function SimplePieChart({ Title, PieData }: { Title: string, PieData: PieRow[] }) {

    const first = useRef(true);
    useEffect(() => { first.current = false; }, []);
    const RADIAN = Math.PI / 180;

    const renderInnerLabel = (props: PieLabelRenderProps) => {
        const { cx, cy, midAngle, innerRadius, outerRadius, percent, name } = props;
        if (percent == null || percent === 0) return <></>;

        const radius = Number(innerRadius) + (Number(outerRadius) - Number(innerRadius)) * 0.6;
        const x = Number(cx) + radius * Math.cos(-(midAngle ?? 0) * RADIAN);
        const y = Number(cy) + radius * Math.sin(-(midAngle ?? 0) * RADIAN);

        return (
            <text
                x={x}
                y={y}
                fill="rgb(22, 22, 22)"
                textAnchor="middle"
                dominantBaseline="central"
                fontSize={13}
                fontWeight={700}
            >
                {`${name} ${(percent * 100).toFixed(0)}%`}
            </text>
        );
    };
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
                        stroke="rgb(124, 139, 151)"
                        strokeWidth={2}
                        label={renderInnerLabel}
                        labelLine={false}
                        // shape={renderShape}
                        animationDuration={300}
                        isAnimationActive={ false}
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
        </div >
    );
}

