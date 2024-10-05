"use client";

import { useAuth } from '@/context/AuthContext';
import { useRouter } from 'next/navigation';
import { useEffect, useContext } from 'react';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card';
import { BarChart, Bar, CartesianGrid, XAxis } from 'recharts';
import * as React from "react";
import { Label, Pie, PieChart } from "recharts";
import { ChartContainer, ChartTooltip, ChartTooltipContent } from "@/components/ui/chart";
import { GlobalContext } from '@/context/GlobalContext';

const integrationsData = [
    {
      name: "Temu LAX",
      createdAt: "2023-01-01",
      lastUpdate: "2023-02-01",
      totalTonnage: 1254,
      barChartData: [
        { month: 'Jan', shipments: 20000, tonnage: 223 },
        { month: 'Feb', shipments: 25000, tonnage: 230 },
        { month: 'Mar', shipments: 30000, tonnage: 100 },
        { month: 'Apr', shipments: 35000, tonnage: 220 },
        { month: 'May', shipments: 40000, tonnage: 22 },
        { month: 'Jun', shipments: 45000, tonnage: 22},
      ],
      pieChartData: [
        { type: "cleared", shipments: 27500, fill: "var(--color-cleared)" },
        { type: "blocked", shipments: 200, fill: "var(--color-blocked)" },
      ],
    },
    {
      name: "Shaoke AMS",
      createdAt: "2023-03-01",
      lastUpdate: "2023-04-01",
      totalShipments: 200000,
      totalTonnage: 9343,
      barChartData: [
        { month: 'Jan', shipments: 50000, tonnage: 250000 },
        { month: 'Feb', shipments: 60000, tonnage: 300000 },
        { month: 'Mar', shipments: 70000, tonnage: 350000 },
      ],
      pieChartData: [
        { type: "cleared", shipments: 89000, fill: "var(--color-cleared)" },
        { type: "blocked", shipments: 20, fill: "var(--color-blocked)" },
      ],
    },
    {
      name: "Shein IST",
      createdAt: "2023-05-01",
      lastUpdate: "2023-06-01",
      totalTonnage: 7532,
      barChartData: [
        { month: 'Jan', shipments: 40000, tonnage: 200000 },
        { month: 'Feb', shipments: 45000, tonnage: 220000 },
        { month: 'Mar', shipments: 50000, tonnage: 250000 },
      ],
      pieChartData: [
        { type: "cleared", shipments: 55000, fill: "var(--color-cleared)" },
        { type: "blocked", shipments: 200, fill: "var(--color-blocked)" },
      ],
    },
    // {
    //   name: "Temu FRA",
    //   createdAt: "2024-05-01",
    //   lastUpdate: "2024-10-05",
    //   totalTonnage: 2324,
    //   barChartData: [
    //     { month: 'May', shipments: 20000, tonnage: 100000 },
    //     { month: 'Aug', shipments: 25000, tonnage: 120000 },
    //     { month: 'Oct', shipments: 30000, tonnage: 150000 },
    //   ],
    //   pieChartData: [
    //     { type: "cleared", shipments: 27225, fill: "var(--color-cleared)" },
    //     { type: "blocked", shipments: 200, fill: "var(--color-blocked)" },
    //   ],
    // },
  ];

const chartConfig = {
  shipments: {
    label: "shipments",
  },
  cleared: {
    label: "Cleared Shipments",
    color: "#36b47E",
  },
  blocked: {
    label: "blocked Shipments",
    color: "#D35178",
  },
};

function PieChartComponent(data) {
  const totalshipments = React.useMemo(() => {
    return data.data.reduce((acc, curr) => acc + curr.shipments, 0);
  }, []);

  return (
    <div style={{ width: '100%', height: '100%', minWidth: '150px', minHeight: '150px' }}>
      <ChartContainer
        config={chartConfig}
        className="mx-auto aspect-square max-h-[250px]"
      >
        <PieChart>
          <ChartTooltip
            cursor={false}
            content={<ChartTooltipContent hideLabel />}
          />
          <Pie
            data={data.data}
            dataKey="shipments"
            nameKey="type"
            innerRadius={30}
            outerRadius={35}
            fill="#8884d8"
            paddingAngle={5}
            strokeWidth={5}
          >
            <Label
              content={({ viewBox }) => {
                if (viewBox && "cx" in viewBox && "cy" in viewBox) {
                  return (
                    <text
                      x={viewBox.cx}
                      y={viewBox.cy}
                      textAnchor="middle"
                      dominantBaseline="middle"
                       fill="#f0f0f0"
                    >
                      {totalshipments}
                    </text>
                  );
                }
                return null;
              }}
            />
          </Pie>
        </PieChart>
      </ChartContainer>
    </div>
  );
}
export default function Integrations() {
  const { isLoggedIn } = useAuth();
  const router = useRouter();
  const { username } = useContext(GlobalContext);


  const handleAddNewIntegrationClick = () => {
    router.push('/integrations/create_integration');
  };

  const routeToIntegration = () => {
    router.push('/integrations/integration');
  }

  useEffect(() => {
    if (!isLoggedIn) {
      router.push('/login');
    }
  }, [isLoggedIn, router]);

  if (!isLoggedIn) {
    return <p>Redirecting to login...</p>;
  }

  return (
    <div className="flex justify-center items-center min-h-screen">
      <div className="flex flex-col items-center gap-5">
        <h1 className="mt-16 mb-8 text-4xl">Welcome {username}</h1>
        <div className="container p-5 rounded-lg inline-block rounded-lg">
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
            {integrationsData.map((integration, index) => (
              <Card key={index} className="w-[370px] h-[370px] cursor-pointer hover:shadow-lg transition-shadow duration-300 rounded-lg" onClick={routeToIntegration}>
                <CardHeader>
                  <CardTitle className="text-2xl">{integration.name}</CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="w-full">
                    <p><strong>active since:</strong> {integration.createdAt}</p>
                    <p><strong>last update:</strong> {integration.lastUpdate}</p>
                  </div>
                  <div className="flex justify-center gap-5 mt-5">
                    <div className="h-38 w-38">
                      <BarChart width={150} height={150} data={integration.barChartData}>
                        <CartesianGrid vertical={false} />
                        <XAxis dataKey="month" tickLine={false} tickMargin={10} axisLine={false} tickFormatter={(value) => value.slice(0, 3)} />
                        <Bar dataKey="shipments" fill="#A08FDB" radius={4} />
                      </BarChart>
                    </div>
                    <div className="h-38 w-38">
                      <PieChartComponent data={integration.pieChartData} />
                    </div>
                  </div>
                  <div className="mt-5 text-center">
                    <p><strong>Total Tonnage:</strong> {integration.totalTonnage}</p>
                  </div>
                </CardContent>
              </Card>
            ))}
            <Card className="w-[370px] h-[370px] cursor-pointer hover:shadow-lg transition-shadow duration-300 flex justify-center items-center rounded-lg" onClick={handleAddNewIntegrationClick}>
              <CardContent>
                <div className="text-center">
                  <p className="text-4xl">+</p>
                  <p>Add New Integration</p>
                </div>
              </CardContent>
            </Card>
          </div>
        </div>
      </div>
    </div>
  );
}