BEGIN;

DROP TABLE IF EXISTS messages;

CREATE TABLE messages (
    id TEXT PRIMARY KEY,
    phone_number TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    sent_at TIMESTAMP,
    status TEXT NOT NULL DEFAULT 'pending'
);

INSERT INTO messages (id, phone_number, content, created_at, sent_at, status) VALUES
    ('1', '+905551111111', 'Merhaba, bu bir test mesajıdır.',  NOW(), NULL, 'pending'),
    ('2', '+905552222222', 'Toplantınız saat 15:00’te.', NOW() - INTERVAL '2 days', NOW() - INTERVAL '1 day 3 hours', 'pending'),
    ( '3', '+905553333333', 'Şifrenizi sıfırlamak için bağlantıya tıklayın.', NOW() - INTERVAL '1 day', NOW() - INTERVAL '23 hours', 'pending'),
    ('4', '+905554444444', 'Kargonuz yola çıktı.',  NOW(), NULL, 'pending'),
    ('5', '+905555555555', 'Bugün hava yağmurlu olacak.', NOW() - INTERVAL '4 days', NOW() - INTERVAL '3 days 22 hours', 'pending'),
    ('6', '+905556666666', 'Yeni kampanyamız başladı!',  NOW(), NULL, 'pending'),
    ('7', '+905557777777', 'Randevunuz onaylandı.', NOW() - INTERVAL '3 days', NOW() - INTERVAL '2 days 20 hours', 'pending'),
    ('8', '+905558888888', 'Lütfen uygulamamızı değerlendirin.',  NOW(), NULL, 'pending'),
    ('9', '+905559999999', 'Ödemeniz başarıyla alındı.', NOW() - INTERVAL '6 days', NOW() - INTERVAL '5 days 12 hours', 'pending'),
    ('10', '+905550000000', 'Bu mesaj sistem tarafından oluşturuldu.',  NOW(), NULL, 'pending'),
    ('11', '+905560000001', 'Sistem güncellemesi başarıyla tamamlandı.', NOW() - INTERVAL '5 days', NOW() - INTERVAL '4 days 23 hours', 'pending'),
    ('12', '+905560000002', 'Yeni şifre talebiniz alındı.', NOW() - INTERVAL '4 days', NOW() - INTERVAL '3 days 20 hours', 'pending'),
    ('13', '+905560000003', 'Kampanya sadece bugün geçerlidir!',  NOW() - INTERVAL '1 days', NULL, 'pending'),
    ('14', '+905560000004', 'Yarın görüşmek üzere!',  NOW(), NULL, 'pending'),
    ('15', '+905560000005', 'Üyelik başarıyla oluşturuldu.', NOW() - INTERVAL '2 days', NOW() - INTERVAL '2 days 1 hours', 'pending'),
    ('16', '+905560000006', 'Destek ekibimiz size ulaşacak.',  NOW(), NULL, 'pending'),
    ('17', '+905560000007', 'Lütfen bilgilerinizi kontrol edin.', NOW() - INTERVAL '7 days', NOW() - INTERVAL '6 days 22 hours', 'pending'),
    ('18', '+905560000008', 'Tek kullanımlık şifreniz: 843729', NOW() - INTERVAL '30 minutes', NOW() - INTERVAL '28 minutes', 'pending'),
    ('19', '+905560000009', 'Faturanız e-posta adresinize gönderildi.', NOW() - INTERVAL '3 days', NOW() - INTERVAL '2 days 18 hours', 'pending'),
    ('20', '+905560000010', 'Etkinliğimiz başladı, bekliyoruz!',  NOW(), NULL, 'pending'),
    ('21', '+905560000011', 'Yeni yorumunuz onaylandı.', NOW() - INTERVAL '5 days', NOW() - INTERVAL '4 days 15 hours', 'pending'),
    ('22', '+905560000012', 'Bugün size özel %20 indirim!',  NOW(), NULL, 'pending'),
    ('23', '+905560000013', 'Randevunuz iptal edildi.', NOW() - INTERVAL '2 days', NOW() - INTERVAL '2 days', 'pending'),
    ('24', '+905560000014', 'Şirket içi bilgilendirme: Toplantı saat 11:00',  NOW(), NULL, 'pending'),
    ('25', '+905560000015', 'Anketimize katıldığınız için teşekkürler.', NOW() - INTERVAL '1 days', NOW() - INTERVAL '1 days', 'pending'),
    ('26', '+905560000016', 'Aboneliğiniz yenilendi.', NOW() - INTERVAL '6 days', NOW() - INTERVAL '5 days 12 hours', 'pending'),
    ('27', '+905560000017', 'İade işleminiz başlatıldı.', NOW() - INTERVAL '3 days', NOW() - INTERVAL '2 days 6 hours', 'pending'),
    ('28', '+905560000018', 'Son giriş saatiniz: 09:23', NOW() - INTERVAL '1 days', NOW() - INTERVAL '20 hours', 'pending'),
    ('29', '+905560000019', 'Yeni cihazdan giriş yapıldı.',  NOW(), NULL, 'pending'),
    ('30', '+905560000020', 'Gönderiniz teslim edildi.', NOW() - INTERVAL '2 days', NOW() - INTERVAL '2 days 2 hours', 'pending'),
    ('31', '+905560000021', 'Ücretsiz deneme süreniz sona eriyor.',  NOW(), NULL, 'pending'),
    ('32', '+905560000022', 'Adres bilgileriniz başarıyla güncellendi.', NOW() - INTERVAL '1 days', NOW() - INTERVAL '22 hours', 'pending'),
    ('33', '+905560000023', 'Yeni bir mesajınız var.', NOW() - INTERVAL '5 hours', NOW() - INTERVAL '4 hours 45 minutes', 'pending'),
    ('34', '+905560000024', 'Soru-cevap etkinliğimiz başladı!',  NOW(), NULL, 'pending'),
    ('35', '+905560000025', 'Giriş yapmayı unuttunuz mu?',  NOW(), NULL, 'pending'),
    ('36', '+905560000026', 'Ürün stoğa geri geldi.', NOW() - INTERVAL '3 days', NOW() - INTERVAL '2 days 12 hours', 'pending'),
    ('37', '+905560000027', 'İndiriminiz sepetinize uygulandı.', NOW() - INTERVAL '1 days', NOW() - INTERVAL '21 hours', 'pending'),
    ('38', '+905560000028', 'Kullanım koşullarımız güncellendi.',  NOW(), NULL, 'pending'),
    ('39', '+905560000029', 'Eğitim kaydınız başarıyla alındı.', NOW() - INTERVAL '4 days', NOW() - INTERVAL '3 days 10 hours', 'pending'),
    ('40', '+905560000030', 'Hizmet puanlamasını unutmayın!',  NOW(), NULL, 'pending');

COMMIT;