import { createBrowserRouter } from "react-router-dom";
import Root from "../components/Root.tsx";
import HomePage from "../pages/home/index.tsx";
import LoginPage from "../pages/login/index.tsx";
import SignupPage from "../pages/signup/index.tsx";

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
    ],
  },
]);
