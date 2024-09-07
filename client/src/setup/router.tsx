import { createBrowserRouter } from "react-router-dom";

import Root from "../components/Root.tsx";
import ProtectedRoute from "../components/navigation/ProtectedRoute.tsx";

import { ConfirmEmailPage, LoginPage, SignupPage } from "../pages/auth";
import DashboardPage from "../pages/dashboard";
import HomePage from "../pages/home";

export const router = createBrowserRouter([
  {
    path: "/",
    element: <Root />,
    children: [
      {
        path: "/",
        element: <HomePage />,
      },
      {
        path: "/login",
        element: <LoginPage />,
      },
      {
        path: "/signup",
        element: <SignupPage />,
      },
      {
        path: "/dashboard",
        element: <ProtectedRoute />,
        children: [
          {
            path: "/dashboard",
            element: <DashboardPage/>,
          },
        ],
      },
      {
        path: "/auth",
        children: [
          {
            path: "confirm-email",
            element: <ConfirmEmailPage />,
          },
        ],
      },
    ],
  },
]);
