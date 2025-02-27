import DyoHead from '@app/components/main/dyo-head'
import { TeamRoutesProvider } from '@app/providers/team-routes'
import { WebSocketProvider } from '@app/providers/websocket'
import '@app/styles/global.css'
import { AppProps } from 'next/app'
import { Toaster } from 'react-hot-toast'

const CruxApp = ({ Component, pageProps }: AppProps) => (
  <>
    <DyoHead />
    <WebSocketProvider>
      <TeamRoutesProvider pageProps={pageProps}>
        <Toaster
          toastOptions={{
            error: {
              icon: null,
              className: '!bg-error-red',
              style: {
                color: 'white',
              },
              position: 'top-center',
            },
          }}
        />

        <Component {...pageProps} />
      </TeamRoutesProvider>
    </WebSocketProvider>
  </>
)

export default CruxApp
