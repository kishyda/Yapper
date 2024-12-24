import styles from './styles.module.css'

export default function topBar() {
    return (
        <div className={styles.topBar}>
            <div className={styles.logo}>LOGO</div>
            <div className={styles.menu}>
                <div>Account</div>
                <div>Signin</div>
                <div>Settings</div>
                <div>Home</div>
            </div>
        </div >
    );
}