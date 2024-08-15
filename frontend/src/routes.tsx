import { createBrowserRouter } from "react-router-dom";
import { DashboardLayout } from "@/layout/Layout";
import { Starmap } from "@/pages/starmap/Starmap";
import { Home } from "@/pages/home/Home";
import { Ships } from "@/pages/ships/Ships";

export const router = createBrowserRouter([
  {
    path: "/",
    element: <DashboardLayout></DashboardLayout>,
    children: [
      {
        path: "/starmap",
        element: <Starmap />,
      },
      {
        path: "/home",
        element: <Home />,
      },
      {
        path: "/ships",
        element: <Ships />,
      },
    ],
  },
]);
