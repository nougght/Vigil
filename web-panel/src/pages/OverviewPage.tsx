import { useEffect, useState } from "react"
import { Link } from "react-router-dom";
import { useOverview } from "../hooks/useOverview";
import { convertBytesToGB } from "@monitoring-system/shared/src/util/units";
import SimplePieChart from "../components/pieChart";
import { CountDistToPieList } from "../domain/overview";

export const OverviewPage = () => {
    const [warning, setWarning] = useState<string | null>()
    const [info, _setInfo] = useState<string | null>()
    const {
        data: overview,
        isPending: isOverviewLoading,
        isError: _isOverviewError,
        error: _overviewError,
        isFetching: _isOverviewFetching,
    } = useOverview();


    useEffect(() => {
        if (overview?.error != null) {
            setWarning(`ошибка:${overview?.error?.status} ${overview?.error?.message}`)
        }
        // if (overview?.overview != null) {
        //     setInfo("Данные успешно загружены")
        // }
    }, [overview]);

    if (isOverviewLoading) {
        return <div>Загрузка...</div>;
    }

    return (

        <div>
            <div>
                <h1>Обзор</h1>
                {overview?.overview?.summary != null && (
                    <div>
                        <div>
                            <label >Общее количество агентов </label>
                            <span>{overview?.overview?.summary?.totalAgents}</span>
                            <br />
                            <label>Онлайн </label>
                            <span>{overview?.overview?.summary?.onlineAgents}</span>
                            <br />
                            <label>Среднее использование CPU </label>
                            <span>{overview?.overview?.summary?.averageCPUUsage.toFixed(2)}%</span>
                            <br />
                            <label>Среднее использование памяти </label>
                            <span>{convertBytesToGB(overview?.overview?.summary?.averageMemoryUsage).toFixed(2)} GB</span>
                        </div>
                        <div className="overview-distributions">
                            <label>Распределение использования CPU</label>
                            <br />
                            <label>Высокое({overview?.overview?.summary?.cpuUsageDistribution?.high?.percent}%) </label>
                            <span>{overview?.overview?.summary?.cpuUsageDistribution?.high?.count}</span>
                            <br />
                            <label>Среднее({overview?.overview?.summary?.cpuUsageDistribution?.medium?.percent}%) </label>
                            <span>{overview?.overview?.summary?.cpuUsageDistribution?.medium?.count}</span>
                            <br />
                            <label>Низкое({overview?.overview?.summary?.cpuUsageDistribution?.low?.percent}%) </label>
                            <span>{overview?.overview?.summary?.cpuUsageDistribution?.low?.count}</span>
                            <br />
                            <label>Распределение использования памяти</label>
                            <br />
                            <label>Высокое({overview?.overview?.summary?.memoryUsageDistribution?.high?.percent}%) </label>
                            <span>{overview?.overview?.summary?.memoryUsageDistribution?.high?.count}</span>
                            <br />
                            <label>Среднее({overview?.overview?.summary?.memoryUsageDistribution?.medium?.percent}%) </label>
                            <span>{overview?.overview?.summary?.memoryUsageDistribution?.medium?.count}</span>
                            <br />
                            <label>Низкое({overview?.overview?.summary?.memoryUsageDistribution?.low?.percent}%) </label>
                            <span>{overview?.overview?.summary?.memoryUsageDistribution?.low?.count}</span>
                            <br />
                            <label>Распределение использования диска</label>
                            <br />
                            <label>Высокое({overview?.overview?.summary?.diskUsageDistribution?.high?.percent}%) </label>
                            <span>{overview?.overview?.summary?.diskUsageDistribution?.high?.count}</span>
                            <br />
                            <label>Среднее({overview?.overview?.summary?.diskUsageDistribution?.medium?.percent}%) </label>
                            <span>{overview?.overview?.summary?.diskUsageDistribution?.medium?.count}</span>
                            <br />
                            <label>Низкое({overview?.overview?.summary?.diskUsageDistribution?.low?.percent}%) </label>
                            <span>{overview?.overview?.summary?.diskUsageDistribution?.low?.count}</span>

                            <div className="dist-pies">

                                <SimplePieChart
                                Title="Использование CPU"
                                PieData={CountDistToPieList(overview.overview.summary.cpuUsageDistribution)}
                                />
                                
                                <SimplePieChart
                                Title="Использование памяти"
                                PieData={CountDistToPieList(overview.overview.summary.memoryUsageDistribution)}
                                />
                                
                                <SimplePieChart
                                Title="Использование диска"
                                PieData={CountDistToPieList(overview.overview.summary.diskUsageDistribution)}
                                />
                            </div>
                        </div>

                        <div className="topN">
                            <div className="topN-cpu">
                                <label>Топ по использованию CPU</label>
                                <br />
                                <table>
                                    {/* <thead>
                                        <tr>
                                            <th>Агент</th>
                                            <th>Использование CPU</th>
                                        </tr>
                                    </thead> */}
                                    <tbody>
                                        {overview?.overview?.topN?.cpuUsage?.map((agent) => (
                                            agent.name != "" &&
                                            <tr key={agent.id}>
                                                <td><Link to={`/agents/${agent.id}`}>{agent.name}</Link></td>
                                                <td>{agent.cpuUsage.toFixed(2)}%</td>
                                            </tr>
                                        ))}
                                    </tbody>
                                </table>
                            </div>
                            <div className="topN-memory">
                                <label>Топ по использованию памяти</label>
                                <br />
                                <table>
                                    {/* <thead>
                                        <tr>
                                            <th>Агент</th>
                                            <th>Использование памяти</th>
                                        </tr>
                                    </thead> */}
                                    <tbody>
                                        {overview?.overview?.topN?.memoryUsage?.map((agent) => (
                                            agent.name != "" &&
                                            <tr key={agent.id}>
                                                <td><Link to={`/agents/${agent.id}`}>{agent.name}</Link></td>
                                                <td>{convertBytesToGB(agent.memoryUsage).toFixed(2)} GB</td>
                                            </tr>
                                        ))}
                                    </tbody>
                                </table>
                            </div>
                        </div>
                    </div>
                )}
            </div>

            {
                warning != null &&
                <div>
                    <p>{warning}</p>
                </div>
            }
            {
                info != null &&
                <div className="infoMessage">
                    <p>{info}</p>
                </div>
            }
        </div>
    )
}