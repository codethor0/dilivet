import { BrowserRouter, Routes, Route, Link } from 'react-router-dom'
import Dashboard from './pages/Dashboard'
import Verify from './pages/Verify'
import KATVerify from './pages/KATVerify'
import './App.css'

function App() {
  return (
    <BrowserRouter>
      <div className="app">
        <header className="app-header">
          <div className="header-content">
            <Link to="/" className="logo-link">
              <img src="/dilivet-logo.png" alt="DiliVet Logo" className="logo" />
              <h1>DiliVet</h1>
            </Link>
            <nav className="nav">
              <Link to="/" className="nav-link">Dashboard</Link>
              <Link to="/verify" className="nav-link">Verify Signature</Link>
              <Link to="/kat-verify" className="nav-link">KAT Verification</Link>
            </nav>
          </div>
        </header>
        <main className="app-main">
          <Routes>
            <Route path="/" element={<Dashboard />} />
            <Route path="/verify" element={<Verify />} />
            <Route path="/kat-verify" element={<KATVerify />} />
          </Routes>
        </main>
      </div>
    </BrowserRouter>
  )
}

export default App

