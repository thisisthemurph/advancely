import { createBrowserRouter } from "react-router-dom";

import Root from "../components/Root.tsx";
import HomePage from "../pages/home/index.tsx";
import LoginPage from "../pages/login/index.tsx";
import SignupPage from "../pages/signup/index.tsx";
import DashboardPage from "../pages/dashboard/index.tsx";
import ConfirmEmailPage from "../pages/auth/ConfirmEmailPage.tsx";

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
        element: <DashboardPage />,
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
