import { Link } from "react-router-dom";
import MobileNavigation from "./Navigation";
import Logo from "./ui/Logo";

function Header() {
  return (
    <header className="flex justify-between shadow-lg p-8">
      <Link to="/" className="text-2xl">
        <Logo size="sm" />
      </Link>
      <MobileNavigation />
    </header>
  );
}

export default Header;
