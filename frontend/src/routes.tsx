import { DashboardLayout } from "@/layout/Layout";
import { Home } from "@/pages/home/Home";
import { Ships } from "@/pages/ships/Ships";
import { Starmap } from "@/pages/starmap/Starmap";
import {
  createBrowserRouter,
  Navigate,
  Outlet,
  Params,
} from "react-router-dom";
import { BreadcrumbData } from "./hooks/use-breadcrumbs";
import { Wallet } from "./pages/wallet/Wallet";

export const router = createBrowserRouter([
  {
    path: "/",
    element: <DashboardLayout></DashboardLayout>,
    children: [
      {
        path: "/",
        element: <Navigate to={"/home"} />,
      },
      {
        path: "/starmap",
        element: <Starmap />,
        handle: {
          crumb: () => ({ name: "Starmap", link: "starmap" }) as BreadcrumbData,
        },
      },
      {
        path: "/home",
        element: <Home />,
        handle: {
          crumb: () => ({ name: "Home", link: "home" }) as BreadcrumbData,
        },
      },
      {
        path: "/ships",
        element: <Ships />,
        handle: {
          crumb: () => ({ name: "Ships", link: "ships" }) as BreadcrumbData,
        },
      },
      {
        path: "wallet",
        element: <Outlet />,
        handle: {
          crumb: () => ({ name: "Wallet", link: "wallet" }) as BreadcrumbData,
        },
        children: [
          {
            path: "",
            element: <Navigate to="transactions" />,
          },
          {
            path: "transactions",
            element: <Wallet />,
          },
          {
            path: "transactions/:id",
            element: <Wallet />,
            handle: {
              crumb: (params: Params) => ({
                name: params.id,
                link: "wallet/transactions/" + params.id,
              }),
            },
          },
        ],
      },
    ],
  },
]);
