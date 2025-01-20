import { HashRouter, Route, Routes } from "react-router-dom"
import { INDEX_PAGE_PATH, IndexPage } from "./pages/index/index_page"
import { HABIT_MANAGER_PAGE_PATH, HabitManagerPage } from "./pages/habit_manager/habit_manager_page"
import { AppLayout } from "./components/app_layout"

function App() {
    return <HashRouter basename={"/"}>
    <Routes>
        <Route path="/" element={<AppLayout />}>
            <Route path={INDEX_PAGE_PATH}  element={<IndexPage/>} />
            <Route path={HABIT_MANAGER_PAGE_PATH} element={<HabitManagerPage/>} />
        </Route>
    </Routes>
  </HashRouter>
}

export default App
