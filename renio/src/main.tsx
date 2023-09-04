import React from 'react'
import ReactDOM from 'react-dom/client'
import App from './App.tsx'
import './conf.ts'
import './assets/main.less'
import './assets/globals.less'

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
)