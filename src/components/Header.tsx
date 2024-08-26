import { Link } from "react-router-dom";
import Navigation from "./Navigation";
import Logo from "./ui/Logo";

function Header() {
  return (
    <header className="flex justify-between items-center shadow-lg p-4">
      <Link to="/" className="text-2xl">
        <Logo size="sm" />
      </Link>
      <Navigation />
    </header>
  );
}

export default Header;
