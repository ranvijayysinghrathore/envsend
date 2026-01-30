import { ImageResponse } from 'next/og';
 
export const runtime = 'edge';
 
export const alt = 'EnvSend - Zero-Knowledge Secret Transfer';
export const size = {
  width: 1200,
  height: 630,
};
 
export const contentType = 'image/png';
 
export default async function Image() {
  return new ImageResponse(
    (
      <div
        style={{
          background: '#050505',
          width: '100%',
          height: '100%',
          display: 'flex',
          flexDirection: 'column',
          alignItems: 'center',
          justifyContent: 'center',
          fontFamily: 'sans-serif',
          position: 'relative',
        }}
      >
        {/* Ambient Gradient Background */}
        <div
          style={{
            position: 'absolute',
            top: '-50%',
            left: '50%',
            transform: 'translate(-50%, 0)',
            width: '1000px',
            height: '1000px',
            background: 'radial-gradient(circle, rgba(255,85,0,0.15) 0%, rgba(0,0,0,0) 70%)',
            filter: 'blur(80px)',
          }}
        />

        {/* Content Container */}
        <div
          style={{
            display: 'flex',
            flexDirection: 'column',
            alignItems: 'center',
            justifyContent: 'center',
            zIndex: 10,
          }}
        >
          {/* Logo */}
          <div
            style={{
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              marginBottom: 40,
            }}
          >
             <img
              src="https://envsend.vercel.app/icon.png"
              alt="EnvSend Logo"
              width={60}
              height={60}
              style={{
                borderRadius: 12,
                marginRight: 20,
              }}
            />
            <div
              style={{
                fontSize: 60,
                fontWeight: 800,
                color: 'white',
                letterSpacing: '-2px',
              }}
            >
              EnvSend
            </div>
          </div>

          <div
            style={{
              fontSize: 32,
              color: 'rgba(255,255,255,0.6)',
              textAlign: 'center',
              maxWidth: '800px',
              fontWeight: 400,
              lineHeight: 1.4,
            }}
          >
            Zero-Knowledge Secret Transfer
          </div>
          
           <div
            style={{
              marginTop: 20,
              fontSize: 24,
              color: 'rgba(255,255,255,0.4)',
              textAlign: 'center',
              maxWidth: '600px',
              fontWeight: 400,
            }}
          >
            No servers. No logs. Pure AES-256.
          </div>
        </div>

        {/* Bottom Decoration */}
        <div
          style={{
            position: 'absolute',
            bottom: 40,
            display: 'flex',
            alignItems: 'center',
            gap: 10,
          }}
        >
           <div style={{ width: 8, height: 8, borderRadius: '50%', background: '#22c55e' }} />
           <div style={{ color: '#22c55e', fontSize: 16, fontWeight: 600 }}>End-to-End Encrypted</div>
        </div>
      </div>
    ),
    {
      ...size,
    }
  );
}
