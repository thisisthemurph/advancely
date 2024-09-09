import { useNavigate } from "react-router-dom";
import { useMediaQuery } from "usehooks-ts";

import DesktopNavigation from "./DesktopNavigation";
import MobileNavigation from "./MobileNavigation";
import { useAuth } from "../../hooks/useAuth.tsx";
import { navigationMenuItems, NavLinks } from "./navigationMenuItems.ts";

export interface NavigationAuthenticationProps {
  isAuthenticated: boolean;
  logout: () => Promise<void>;
}

export interface NavigationProps extends NavigationAuthenticationProps {
  menuItems: NavLinks;
}

function Navigation() {
  const navigate = useNavigate();
  const { isAuthenticated, logout } = useAuth();
  const isDesktop = useMediaQuery("(min-width: 640px)");

  const handleLogout = async () => {
    logout().finally(() => { navigate("/login"); });
  }

  const filteredMenu = navigationMenuItems.filter((link) => {
    switch (link.showWhen) {
      case "always":
      case undefined:
        return true;
      case "authenticated":
        return isAuthenticated;
      case "unauthenticated":
        return !isAuthenticated;
      default:
        return false;
    }
  });

  return isDesktop ? (
    <DesktopNavigation menuItems={filteredMenu} isAuthenticated={isAuthenticated} logout={handleLogout} />
  ) : (
    <MobileNavigation menuItems={filteredMenu} isAuthenticated={isAuthenticated} logout={handleLogout} />
  );
}

export default Navigation;
