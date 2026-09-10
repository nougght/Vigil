import { Legend, Pie, PieChart, ResponsiveContainer, Sector, Tooltip } from "recharts";


interface PieRow{
    Name: string,
    Value: number,
    fill: string,
}

export default function SimplePieChart({Title, PieData}: {Title: string, PieData: PieRow[]}) {
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
            label={({ name, percent }) => `${name} ${((percent ?? 0) * 100).toFixed(0)}%`}
            shape={(props) => {
              const { fill, payload, ...rest } = props;
              // fallback to payload color if available, otherwise use standard fill
              return <Sector {...rest} fill={payload?.fill || fill} />;
            }}
          />
          <Tooltip />
          <Legend />
        </PieChart>
      </ResponsiveContainer>
    </div>
  );
}

